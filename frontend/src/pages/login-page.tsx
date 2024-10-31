import axios, { AxiosError } from "axios";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { useToast } from "@/hooks/use-toast";
import { useQuery } from "@tanstack/react-query";
import { LoaderCircle } from "lucide-react";

const registerFormSchema = z.object({
  email: z
    .string()
    .email({ message: "Ivalid email." })
    .min(1, { message: "This field is required." }),
  password: z.string().min(1, { message: "This field is required." }),
});

const LoginPage = () => {
  const navigate = useNavigate();
  const { toast } = useToast();

  const { isLoading } = useQuery({
    queryKey: ["logged-in-query"],
    queryFn: async () => {
      try {
        const response = await axios.get("/api/user/logged-in");

        if (response.status === 200) {
          navigate("/user-page");
        }

        return true
      } catch (err) {
        toast({
          variant: "destructive",
          title: "Checking whether user is already logged in...",
          description: (err as AxiosError).message,
        });
        return false
      }
    },
  });

  const loginForm = useForm<z.infer<typeof registerFormSchema>>({
    resolver: zodResolver(registerFormSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const registerFormOnSubmit = async (
    values: z.infer<typeof registerFormSchema>
  ) => {
    try {
      const response = await axios.post("/api/user/login", {
        email: values.email,
        password: values.password,
      });

      if (response.status === 200) {
        navigate("/user-page");
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
      <div className="flex">
        <div className="flex items-center justify-center w-1/3 bg-slate-900">
          <p className="text-5xl font-bold text-white">Login Page</p>
        </div>
        <div className="flex items-center justify-center h-screen w-2/3 bg-slate-950">
          <LoaderCircle size={38} className="animate-spin" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex">
      <div className="flex items-center justify-center w-1/3 bg-slate-900">
        <p className="text-5xl font-bold text-white">Login Page</p>
      </div>
      <div className="flex items-center justify-center h-screen w-2/3 bg-slate-950">
        <Form {...loginForm}>
          <form
            onSubmit={loginForm.handleSubmit(registerFormOnSubmit)}
            className="flex flex-col items-center gap-y-4"
          >
            <FormField
              control={loginForm.control}
              name="email"
              render={({ field }) => (
                <FormItem className="transition-all duration-300">
                  <FormControl>
                    <Input
                      type="email"
                      placeholder="Email"
                      autoComplete="off"
                      className="font-semibold transition-all duration-300"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage className="transition-all duration-300" />
                </FormItem>
              )}
            />
            <FormField
              control={loginForm.control}
              name="password"
              render={({ field }) => (
                <FormItem className="transition-all duration-300">
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="Password"
                      autoComplete="off"
                      className="font-semibold transition-all duration-300"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage className="transition-all duration-300" />
                </FormItem>
              )}
            />
            <Button
              variant="secondary"
              type="submit"
              className="w-1/2 font-semibold"
            >
              Login
            </Button>
          </form>
        </Form>
      </div>
    </div>
  );
};

export default LoginPage;
