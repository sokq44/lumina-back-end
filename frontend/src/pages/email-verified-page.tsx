import { useToast } from "@/hooks/use-toast";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { LoaderCircle, MailCheck, MailWarning } from "lucide-react";
import { useParams } from "react-router-dom";

const EmailVerifiedPage = () => {
  const { token } = useParams();
  const { toast } = useToast();

  const { data, isLoading } = useQuery({
    queryKey: ["email-verification"],
    queryFn: async () => {
      if (token) {
        const response = await axios.patch("/api/user/verify-email", {
          token: token,
        });
        if (response.status == 204) return true;
        toast({
          variant: "destructive",
          title: "Error",
          description: `Server responded with status code: ${response.status}`,
        });
        return false;
      }
    },
  });

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-screen bg-slate-900">
        <LoaderCircle className="animate-spin text-white" size={32} />
      </div>
    );
  }

  if (data) {
    return (
      <div className="flex flex-col items-center justify-center h-screen bg-slate-900">
        <span className="flex items-center text-4xl font-bold text-white">
          Your email has been Verified &nbsp; <MailCheck size={32} />
        </span>
        <p className="text-lg text-slate-600">
          Enjoy your journey to knowledge and may the force be with you.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen bg-slate-900">
      <span className="flex items-center text-4xl font-bold text-white">
        Problem with email verification &nbsp; <MailWarning size={32} />
      </span>
      <p className="text-lg text-slate-600">
        We're sorry but it seems that an unexpected problem has occured while
        trying to verify this email address.
      </p>
    </div>
  );
};

export default EmailVerifiedPage;
