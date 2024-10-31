import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";
import { useQuery } from "@tanstack/react-query";
import axios, { AxiosError } from "axios";
import { LoaderCircle } from "lucide-react";
import { useNavigate } from "react-router-dom";

const UserPage = () => {
  const navigate = useNavigate();
  const { toast } = useToast();

  const { data, isLoading } = useQuery({
    queryKey: ["user-info-query"],
    queryFn: async (): Promise<{ [s: string]: string } | undefined> => {
      try {
        const response = await axios.get("/api/user/get-user");

        return {
          username: response.data.username,
          email: response.data.email,
        };
      } catch (err) {
        const error = err as AxiosError;
        if (error.status === 401) navigate("/login");
      }
    },
    retryDelay: 10000,
  });

  const modify = async () => {
    navigate("/modify-user", {
      state: { username: data?.username, email: data?.email },
    });
  };

  const logout = async () => {
    try {
      const response = await axios.delete("/api/user/logout");

      if (response.status == 200) {
        navigate("/login");
      }
    } catch (err) {
      toast({
        variant: "destructive",
        title: "Problem with registering",
        description: (err as AxiosError).message,
      });
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-screen bg-slate-950">
        <LoaderCircle size={38} className="animate-spin" />
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4 items-center justify-center h-screen bg-slate-950">
      <div className="flex items-center justify-center text-white">
        {JSON.stringify(data)}
      </div>
      <div className="flex items-center gap-4">
        <Button variant="secondary" onClick={modify}>
          Modify The Data
        </Button>
        <Button variant="secondary" onClick={logout}>
          Logout
        </Button>
      </div>
    </div>
  );
};

export default UserPage;
