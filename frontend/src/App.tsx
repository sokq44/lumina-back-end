import { createBrowserRouter, RouterProvider } from "react-router-dom";
import EmailVerifiedPage from "./pages/email-verified-page";
import VerifyEmailPage from "@/pages/verify-email-page";
import RegisterPage from "@/pages/register-page";
import { Toaster } from "@/components/ui/toaster";

const router = createBrowserRouter([
  {
    path: "/",
    element: <RegisterPage />,
  },
  {
    path: "/verify-email",
    element: <EmailVerifiedPage />,
  },
  {
    path: "/verify-email-info",
    element: <VerifyEmailPage />,
  },
]);

export default function App() {
  return (
    <>
      <RouterProvider router={router} />
      <Toaster />
    </>
  );
}
