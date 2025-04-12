import { MetaProvider, Title } from "@solidjs/meta";
import { Router } from "@solidjs/router";
import { FileRoutes } from "@solidjs/start/router";
import { Suspense } from "solid-js";
import "./app.css";

export default function App() {
  return (
    <Router
      root={props => (
        <MetaProvider>
          <Title>{"Blue Fox | " + props.location.pathname.replace(/^\//, '').replace(/\//g, ' / ')}</Title>
          <Suspense fallback={<div>Loading...</div>}>{props.children}</Suspense>
        </MetaProvider>
      )}
    >
    <FileRoutes />
    </Router>
  );
}
