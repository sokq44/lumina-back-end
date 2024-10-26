import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";
import axios, { AxiosError } from "axios";
import { useNavigate } from "react-router-dom";

const UserPage = () => {
  const navigate = useNavigate();
  const { toast } = useToast();

  return (
    <div className="flex items-center justify-center h-screen">
      <Button
        onClick={async () => {
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
        }}
      >
        Logout
      </Button>
    </div>
  );
};

export default UserPage;
