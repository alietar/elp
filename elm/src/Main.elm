module Main exposing (main)

import Browser
import Html exposing (Html, div, button, text, p)
import Html.Events exposing (onClick)
import Http
import Carte
import Draw_square
import UserApi
import Interface

-- MODEL

type alias Model =
    { carte : Carte.Model
    , status : String
    , form : Interface.Model
    }


-- INIT

init : () -> ( Model, Cmd Msg)
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
      , status = "Prêt à charger les carrés."
      , form = initialForm
      }
    , carteCmd -- On initialise seulement la carte, pas de carrés au démarrage
    )


-- MSG

type Msg
    = FormMsg Interface.Msg
    | GotSquares (Result Http.Error UserApi.ServerResponse)


-- 4. UPDATE

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        -- Gestion des messages venant de l'Interface (champs texte, bouton valider)
        FormMsg interfaceMsg ->
            case interfaceMsg of
                Interface.Lat val ->
                    let oldForm = model.form
                    in ({ model | form = { oldForm | lat = val, validate = False, typeError = False } }, Cmd.none)

                Interface.Long val ->
                    let oldForm = model.form
                    in ({ model | form = { oldForm | long = val, validate = False, typeError = False } }, Cmd.none)

                Interface.Deniv val ->
                    let oldForm = model.form
                    in ({ model | form = { oldForm | d = val, validate = False, typeError = False } }, Cmd.none)

                Interface.Validate ->
                    -- C'est ICI qu'on fait le lien entre l'Interface et l'API
                    let
                        maybeLat = String.toFloat model.form.lat
                        maybeLng = String.toFloat model.form.long
                        maybeDeniv = String.toFloat model.form.d
                    in
                    case (maybeLat, maybeLng, maybeDeniv) of
                        (Just lat, Just lng, Just deniv) ->
                            let
                                -- Données valides : On prépare la requête
                                apiData : UserApi.StartPointData
                                apiData =
                                    { lat = lat
                                    , lng = lng
                                    , deniv = deniv
                                    }
                                
                                oldForm = model.form
                                newForm = { oldForm | validate = True, typeError = False }
                            in
                            ( { model | status = "Chargement...", form = newForm }
                            , UserApi.fetchSquares apiData GotSquares
                            )

                        _ ->
                            -- Données invalides : On affiche l'erreur
                            let
                                oldForm = model.form
                                newForm = { oldForm | validate = False, typeError = True }
                            in
                            ( { model | form = newForm, status = "Erreur de saisie." }
                            , Cmd.none
                            )

        -- Gestion de la réponse du Serveur
        GotSquares result ->
            case result of
                Ok squares ->
                    let
                        -- Conversion des données (ton code original)
                        toParams sq =
                            { size = sq.size
                            , centerLat = sq.centerLat
                            , centerLng = sq.centerLng
                            }

                        boundsList =
                            List.map (toParams >> Draw_square.computeBounds) squares

                        drawCmds =
                            List.map Carte.drawSquare boundsList

                        zoomCmd = Carte.autoView ()
                        
                        -- Ajout du nettoyage (optionnel si tu l'as implémenté)
                        clearCmd = Carte.clearSquares ()
                    in
                    ( { model | status = "Succès : " ++ String.fromInt (List.length squares) ++ " carrés." }
                    , Cmd.batch (clearCmd :: drawCmds ++ [ zoomCmd ])
                    )

                Err _ ->
                    ( { model | status = "Erreur lors de la récupération des données." }
                    , Cmd.none
                    )


-- 5. VIEW

view : Model -> Html Msg
view model =
    div []
        [ -- Affiche l'interface en convertissant ses messages (FormMsg)
          Html.map FormMsg (Interface.mainView model.form)
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