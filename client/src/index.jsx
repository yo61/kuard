import {createRoot} from 'react-dom/client';
import App from './app'

if (process.env.NODE_ENV !== 'production') {
  console.log('Looks like we are in development mode!');
}

createRoot(document.getElementById("root")).render(<App page={pageContext}/>)
