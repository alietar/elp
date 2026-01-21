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
            ( model, Cmd.none )


        -- Messages venant de la carte (clic)
        MapMsg carteMsg ->
            let
                ( newCarte, carteCmd ) =
                    Carte.update carteMsg model.carte
            in
            case newCarte.clicked of
                Just coord ->
                    let
                        oldForm = model.form
                    in
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
        GotSquares _ ->
            ( model, Cmd.none )




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
