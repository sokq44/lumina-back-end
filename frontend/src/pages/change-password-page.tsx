import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import { zodResolver } from "@hookform/resolvers/zod";
import axios, { AxiosError } from "axios";
import { useForm } from "react-hook-form";
import { z } from "zod";

const changePasswordFormSchema = z
  .object({
    password: z.string().min(1, { message: "This field is required." }),
    repeat: z.string(),
  })
  .refine((data) => data.password === data.repeat, {
    message: "Passwords don't match",
    path: ["repeatPass"],
  });

const ChangePasswordPage = () => {
  const { toast } = useToast();

  const changePasswordForm = useForm<z.infer<typeof changePasswordFormSchema>>({
    resolver: zodResolver(changePasswordFormSchema),
    defaultValues: {
      password: "",
      repeat: "",
    },
  });

  const changePasswordFormOnSubmit = async (
    values: z.infer<typeof changePasswordFormSchema>
  ) => {
    try {
      const response = await axios.patch("/api/user/change-password", {
        password: values.password,
      });
      console.log(response);
    } catch (err) {
      toast({
        variant: "destructive",
        title: "Problem with modifyUsering",
        description: (err as AxiosError).message,
      });
    }
  };

  return (
    <div className="flex flex-col items-center justify-center gap-y-2 h-screen bg-slate-950">
      <Form {...changePasswordForm}>
        <form
          onSubmit={changePasswordForm.handleSubmit(changePasswordFormOnSubmit)}
          className="flex flex-col items-center gap-y-4"
        >
          <FormField
            control={changePasswordForm.control}
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
          <FormField
            control={changePasswordForm.control}
            name="repeat"
            render={({ field }) => (
              <FormItem className="transition-all duration-300">
                <FormControl>
                  <Input
                    type="password"
                    placeholder="Repeat Password"
                    autoComplete="off"
                    className="font-semibold transition-all duration-300"
                    {...field}
                  />
                </FormControl>
                <FormMessage className="transition-all duration-300" />
              </FormItem>
            )}
          />
          <Button variant="secondary" type="submit" className="font-semibold">
            Change Password
          </Button>
        </form>
      </Form>
    </div>
  );
};

export default ChangePasswordPage;
