import React from 'react';
import { render } from 'ink';
import { withFullScreen } from "fullscreen-ink";
import App from './interface/app.js';
withFullScreen(/*#__PURE__*/React.createElement(App, null)).start();