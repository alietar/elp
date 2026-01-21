module Main exposing (main)

import Browser
import Html exposing (Html, div)
import Http

import Carte
import Draw_square
import Interface
import UserApi
import Round



-- MODEL

type alias Model =
    { carte : Carte.Model
    , form : Interface.Model
    , status : String
    }



-- MSG

type Msg
    = FormMsg Interface.Msg
    | MapMsg Carte.Msg
    | GotSquares (Result Http.Error UserApi.ServerResponse)



-- INIT

init : () -> ( Model, Cmd Msg )
init _ =
    let
        ( carteModel, carteCmd ) =
            Carte.init

        initialForm =
            { lat = ""
            , long = ""
            , d = ""
            , validate = False
            , typeError = False
            }
    in
    ( { carte = carteModel
      , form = initialForm
      , status = "Prêt."
      }
    , Cmd.map MapMsg carteCmd
    )



-- UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of

        -- Messages venant de l'interface
        FormMsg interfaceMsg ->
            let
                oldForm =
                    model.form
            in
            case interfaceMsg of

                Interface.Lat val ->
                    ( { model
                        | form =
                            { oldForm
                                | lat = val
                                , validate = False
                                , typeError = False
                            }
                      }
                    , Cmd.none
                    )

                Interface.Long val ->
                    ( { model
                        | form =
                            { oldForm
                                | long = val
                                , validate = False
                                , typeError = False
                            }
                      }
                    , Cmd.none
                    )

                Interface.Deniv val ->
                    ( { model
                        | form =
                            { oldForm
                                | d = val
                                , validate = False
                                , typeError = False
                            }
                      }
                    , Cmd.none
                    )

                Interface.Validate ->
                    let
                        maybeLat =
                            String.toFloat oldForm.lat

                        maybeLon =
                            String.toFloat oldForm.long

                        maybeDeniv =
                            String.toFloat oldForm.d
                    in
                    case ( maybeLat, maybeLon, maybeDeniv ) of
                        ( Just lat, Just lon, Just deniv ) ->
                            let
                                apiData =
                                    { lat = lat
                                    , lng = lon
                                    , deniv = deniv
                                    }
                            in
                            ( { model
                                | status = "Calcul en cours..."
                                , form =
                                    { oldForm
                                        | validate = True
                                        , typeError = False
                                    }
                              }
                            , UserApi.fetchSquares apiData GotSquares
                            )

                        _ ->
                            ( { model
                                | status = "Erreur de saisie"
                                , form =
                                    { oldForm
                                        | typeError = True
                                    }
                              }
                            , Cmd.none
                            )


        -- Messages venant de la carte (clic)
        MapMsg carteMsg ->
            let
                ( newCarte, carteCmd ) =
                    Carte.update carteMsg model.carte

                oldForm =
                    model.form
            in
            case newCarte.clicked of
                Just coord ->
                    ( { model
                        | carte = newCarte
                        , form =
                            { oldForm
                                | lat = Round.round 6 coord.lat
                                , long = Round.round 6 coord.lon
                                , validate = False
                                , typeError = False
                            }
                      }
                    , Cmd.map MapMsg carteCmd
                    )

                Nothing ->
                    ( { model | carte = newCarte }
                    , Cmd.map MapMsg carteCmd
                    )


        -- Réponse de l'API
        GotSquares result ->
            case result of
                Ok squares ->
                    let
                        -- Conversion API → bounds Leaflet
                        boundsList =
                            List.map
                                (\sq ->
                                    Draw_square.computeBounds
                                        { size = sq.size
                                        , centerLat = sq.centerLat
                                        , centerLng = sq.centerLng
                                        }
                                )
                                squares

                        -- Commandes Leaflet
                        drawCmds =
                            List.map Carte.drawSquare boundsList

                        clearCmd =
                            Carte.clearSquares ()

                        zoomCmd =
                            Carte.autoView ()
                    in
                    ( { model
                        | status =
                            "Succès : "
                                ++ String.fromInt (List.length squares)
                                ++ " carrés."
                    }
                    , Cmd.batch (clearCmd :: drawCmds ++ [ zoomCmd ])
                    )

                Err _ ->
                    ( { model | status = "Erreur serveur." }
                    , Cmd.none
                    )




-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.map MapMsg (Carte.subscriptions model.carte)



-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ Html.map FormMsg (Interface.mainView model.form)
        , Carte.view model.carte
        ]



-- MAIN

main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }
