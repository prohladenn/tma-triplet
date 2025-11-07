import { AppRoot } from "@telegram-apps/telegram-ui";
import { miniApp, useLaunchParams, useSignal } from "@tma.js/sdk-react";
import { HashRouter, Navigate, Route, Routes } from "react-router-dom";

import { routes } from "@/navigation/routes.tsx";

export function App() {
  const lp = useLaunchParams();
  const isDark = useSignal(miniApp.isDark);

  return (
    <AppRoot
      appearance={isDark ? "dark" : "light"}
      platform={["macos", "ios"].includes(lp.tgWebAppPlatform) ? "ios" : "base"}
    >
      <HashRouter>
        <Routes>
          {routes.map((route) => (
            <Route key={route.path} {...route} />
          ))}
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
      </HashRouter>
    </AppRoot>
  );
}
