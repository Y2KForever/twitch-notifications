import React from 'react';
import ReactDOM from 'react-dom/client';
import { App } from './App';
import { ThemeProvider } from './components/ThemeProvider';
import { store } from './store/store';
import { Provider } from 'react-redux';
import { CheckConfigFile, CheckAndReadFile } from '../wailsjs/go/main/App';

const rootElement = document.getElementById('root');

const v = CheckAndReadFile();

if (rootElement) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <React.StrictMode>
      <Provider store={store}>
        <div className="fixed-size h-full">
          <ThemeProvider defaultTheme="dark" storageKey="tn-theme">
            <App />
          </ThemeProvider>
        </div>
      </Provider>
    </React.StrictMode>,
  );
}
