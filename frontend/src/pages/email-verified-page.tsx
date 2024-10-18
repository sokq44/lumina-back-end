import { MailCheck } from "lucide-react";

const EmailVerifiedPage = () => {
  return (
    <div className="flex flex-col items-center justify-center h-screen bg-slate-900">
      <span className="flex items-center text-4xl font-bold text-white">
        Your email has been Verified &nbsp; <MailCheck size={32} />
      </span>
      <p className="text-xl text-slate-600">
        Enjoy your journey to knowledge and may the force be with you.
      </p>
    </div>
  );
};

export default EmailVerifiedPage;
