import React from 'react';
import {render} from 'ink';
import { withFullScreen } from "fullscreen-ink";
import App from './interface/app.js';


withFullScreen(<App />).start();