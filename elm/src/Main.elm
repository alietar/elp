module Main exposing (main)

import Browser
import Html exposing (Html, div)
import Http

import Carte
import DrawSquare
import Interface
import UserApi
import Round



-- MODEL

type alias Model =
    { carte : Carte.Model
    , form : Interface.Model
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
            , accuracy = "25m"
            , status = Interface.Idle
            , typeError = False
            }
    in
    ( { carte = carteModel
      , form = initialForm
      }
    , Cmd.map MapMsg carteCmd
    )



-- UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of

        -- ============================
        -- MESSAGES DE L'INTERFACE
        -- ============================

        FormMsg interfaceMsg ->
            let
                oldForm =
                    model.form
            in
            case interfaceMsg of
                Interface.Lat val -> -- mise à jour de latitude si message de l'interface
                    ( { model
                        | form =
                            { oldForm
                                | lat = val
                                , status = Interface.Idle
                                , typeError = False
                            }
                      }
                    , Cmd.none
                    )

                Interface.Long val -> -- mise à jour de longitude si message de l'interface
                    ( { model
                        | form =
                            { oldForm
                                | long = val
                                , status = Interface.Idle
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
                                , status = Interface.Idle
                                , typeError = False
                            }
                      }
                    , Cmd.none
                    )

                Interface.Accuracy val ->
                    ( { model
                        | form =
                            { oldForm
                                | accuracy = val
                                , status = Interface.Idle
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
                        maybeAccuracy =
                            oldForm.accuracy
                                |> String.replace "m" ""
                                |> String.toFloat
                    in
                    case ( (maybeLat, maybeLon), (maybeDeniv, maybeAccuracy) ) of
                        ( (Just lat, Just lon), (Just deniv, Just accuracy) ) ->
                            let
                                apiData =
                                    { lat = lat
                                    , lng = lon
                                    , deniv = deniv
                                    , accuracy = accuracy
                                    }

                                -- 🔹 on déclenche la carte ICI
                                ( newCarte, carteCmd ) =
                                    Carte.update
                                        (Carte.requestMarker -- on fait la requête de marqueur si on remplit des données dans l'interface
                                            { lat = lat
                                            , lon = lon
                                            }
                                        )
                                        model.carte
                            in
                            ( { model
                                | carte = newCarte
                                , form =
                                    { oldForm
                                        | status = Interface.Loading
                                        , typeError = False
                                    }
                              }
                            , Cmd.batch
                                [ UserApi.fetchSquares apiData GotSquares
                                , Cmd.map MapMsg carteCmd
                                ]
                            )

                        _ ->
                            ( { model
                                | form =
                                    { oldForm | typeError = True }
                              }
                            , Cmd.none
                            )


        -- ============================
        -- MESSAGES DE LA CARTE
        -- ============================

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
                                | lat = Round.round 6 coord.lat -- on arrondi à 6 chiffres après la virgule les coordonnées cliquées
                                , long = Round.round 6 coord.lon
                                , status = Interface.Idle
                                , typeError = False
                            }
                      }
                    , Cmd.map MapMsg carteCmd
                    )

                Nothing ->
                    ( { model | carte = newCarte }
                    , Cmd.map MapMsg carteCmd
                    )


        -- ============================
        -- RÉPONSE DE L'API
        -- ============================

        GotSquares result ->
            case result of
                Ok data ->
                    let
                        oldForm = model.form
                        boundsList = -- calcul des coordonnées des 4 sommets des carrées/rectangles
                            List.map
                                (\sq ->
                                    DrawSquare.computeBounds
                                        { size = sq.size * data.tileSize
                                        , centerLat = sq.centerLat
                                        , centerLng = sq.centerLng
                                        }
                                )
                                data.tiles

                        drawCmds =
                            List.map Carte.drawSquare boundsList -- on envoie la commande pour tracer les carrés

                        clearCmd =
                            Carte.clearSquares ()

                        zoomCmd =
                            Carte.autoView ()
                    in
                    ( { model
                        | form = { oldForm | status = Interface.Success }
                      }
                    , Cmd.batch
                        (Cmd.map MapMsg clearCmd
                            :: List.map (Cmd.map MapMsg) drawCmds
                            ++ [ Cmd.map MapMsg zoomCmd ]
                        )
                    )

                Err _ ->
                    let
                        oldForm = model.form
                    in
                    ( { model
                        | form = { oldForm | status = Interface.Error }
                    }
                    , Cmd.none
                    )



-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg -- on écoute les évenements du fichier Carte.elm
subscriptions model =
    Sub.map MapMsg (Carte.subscriptions model.carte)



-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ Html.map FormMsg (Interface.mainView model.form)
        , Html.map MapMsg (Carte.view model.carte)
        ] -- ajout des balises HTML 



-- MAIN

main : Program () Model Msg
main =
    Browser.element
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        }
