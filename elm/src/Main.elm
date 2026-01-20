module Main exposing (main)

import Browser
import Html exposing (Html, div, button, text, p)
import Html.Events exposing (onClick)
import Http
import Carte
import Draw_square
import UserApi


-- MODEL

type alias Model =
    { carte : Carte.Model
    , status : String
    }


-- INIT

init : () -> ( Model, Cmd Msg )
init _ =
    let
        ( carteModel, carteCmd ) =
            Carte.init
    in
    ( { carte = carteModel
      , status = "Prêt à charger les carrés."
      }
    , carteCmd -- On initialise seulement la carte, pas de carrés au démarrage
    )


-- MSG

type Msg
    = ClickFetch
    | GotSquares (Result Http.Error UserApi.ServerResponse)


-- UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        ClickFetch ->
            let
                -- Coordonnées fixes pour l'instant, comme demandé
                fixedData : UserApi.StartPointData
                fixedData =
                    { lat = 45.7838052
                    , lng = 4.871928
                    , deniv = 0.3
                    }
            in
            ( { model | status = "Chargement en cours..." }
            , UserApi.fetchSquares fixedData GotSquares
            )

        GotSquares result ->
            case result of
                Ok squares ->
                    let
                        -- 1. On transforme le type UserApi.Square vers Draw_square.Params
                        toParams : UserApi.Square -> Draw_square.Params
                        toParams sq =
                            { size = sq.size
                            , centerLat = sq.centerLat
                            , centerLng = sq.centerLng
                            }

                        -- 2. On calcule les bornes (Bounds) pour chaque carré
                        boundsList =
                            List.map (toParams >> Draw_square.computeBounds) squares

                        -- 3. On crée une commande pour dessiner chaque carré
                        drawCmds =
                            List.map Carte.drawSquare boundsList

                        zoomCmd = Carte.autoView ()
                        clearCmd = Carte.clearSquares ()

                        
                    in
                    ( { model | status = "Succès : " ++ String.fromInt (List.length squares) ++ " carrés affichés." }
                    , Cmd.batch ( clearCmd :: drawCmds ++ [ zoomCmd ])
                    )

                Err _ ->
                    ( { model | status = "Erreur lors de la récupération des données." }
                    , Cmd.none
                    )


-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ div [  ]
            [ button [ onClick ClickFetch ] [ text "Récupérer et afficher les carrés" ]
            , p [] [ text model.status ]
            ]
        , Carte.view model.carte
        ]


-- MAIN

main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = \_ -> Sub.none
        }