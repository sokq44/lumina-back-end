import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import EmailVerifiedPage from "./pages/email-verified-page";
import VerifyEmailPage from "@/pages/verify-email-page";
import RegisterPage from "@/pages/register-page";
import { Toaster } from "@/components/ui/toaster";
import "./index.css";

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
