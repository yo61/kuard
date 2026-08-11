import React from 'react';

// Lets the page components tell <Disconnected/> that a fetch failed, so it can flash the
// power icon. Replaces the legacy getChildContext/contextTypes pair, which React 19 removed.
// The default is a no-op so a component rendered outside the provider still works.
export default React.createContext({ reportConnError: () => {} });
