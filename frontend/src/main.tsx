import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import EmailVerifiedPage from "./pages/email-verified-page";
import VerifyEmailPage from "@/pages/verify-email-page";
import RegisterPage from "@/pages/register-page";
import { Toaster } from "@/components/ui/toaster";
import "./index.css";
import LoginPage from "./pages/login-page";
import UserPage from "./pages/user-page";

const router = createBrowserRouter([
  {
    path: "/",
    element: <RegisterPage />,
  },
  {
    path: "/verify-email/:token",
    element: <EmailVerifiedPage />,
  },
  {
    path: "/verify-email-info",
    element: <VerifyEmailPage />,
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/user-page",
    element: <UserPage />,
  },
]);

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      refetchOnMount: false,
      refetchOnReconnect: false,
      retry: false,
    },
  },
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
      <Toaster />
    </QueryClientProvider>
  </StrictMode>
);
