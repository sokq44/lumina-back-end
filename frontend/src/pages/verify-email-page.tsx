import { Link, useLocation } from "react-router-dom";

const VerifyEmailPage = () => {
  const { state } = useLocation();
  return (
    <div className="flex">
      <div className="flex items-center justify-center w-1/3 bg-slate-900">
        <p className="text-5xl font-bold text-white">Verify Email Page</p>
      </div>
      <div className="flex items-center justify-center h-screen w-2/3 bg-slate-950">
        <span className="text-white font-semibold w-1/2">
          A verification email has been sent to{" "}
          <span className="text-slate-600">{state.email}</span>. Check your
          inbox for the verification link. If You can't see it, then check the
          spam folder.
          <br />
          <Link className="font-bold underline text-slate-600" to={"/"}>
            Back to the home page
          </Link>
        </span>
      </div>
    </div>
  );
};

export default VerifyEmailPage;
