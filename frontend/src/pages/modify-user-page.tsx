import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useToast } from "@/hooks/use-toast";
import axios, { AxiosError } from "axios";
import { useLocation, useNavigate } from "react-router-dom";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

const modifyUserFormSchema = z.object({
  username: z
    .string()
    .min(1, { message: "Must be at least 5 characters long." })
    .max(50, { message: "Can't be longer than 50 characters." }),
  email: z
    .string()
    .email({ message: "Ivalid email." })
    .min(1, { message: "This field is required." }),
});

const ModifyUserPage = () => {
  const { toast } = useToast();
  const navigate = useNavigate();
  const { state } = useLocation();

  const modifyUserForm = useForm<z.infer<typeof modifyUserFormSchema>>({
    resolver: zodResolver(modifyUserFormSchema),
    defaultValues: {
      username: state.username,
      email: state.email,
    },
  });

  const modifyUserFormOnSubmit = async (
    values: z.infer<typeof modifyUserFormSchema>
  ) => {
    try {
      const response = await axios.patch("/api/user/modify-user", {
        username: values.username,
        email: values.email,
      });

      if (response.status == 200) {
        navigate("/user-page");
      }
    } catch (err) {
      toast({
        variant: "destructive",
        title: "Problem with modifyUsering",
        description: (err as AxiosError).message,
      });
    }
  };

  return (
    <div className="flex flex-col gap-4 items-center justify-center h-screen bg-slate-950">
      <Form {...modifyUserForm}>
        <form
          onSubmit={modifyUserForm.handleSubmit(modifyUserFormOnSubmit)}
          className="flex flex-col items-center gap-y-4"
        >
          <FormField
            control={modifyUserForm.control}
            name="username"
            render={({ field }) => (
              <FormItem className="transition-all duration-300">
                <FormControl>
                  <Input
                    type="text"
                    placeholder="Username"
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
            control={modifyUserForm.control}
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
          <Button
            variant="secondary"
            type="submit"
            className="w-1/2 font-semibold"
          >
            Modify
          </Button>
        </form>
      </Form>
    </div>
  );
};

export default ModifyUserPage;
