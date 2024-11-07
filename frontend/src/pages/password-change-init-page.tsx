import { z } from "zod";
import { useToast } from "@/hooks/use-toast";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Link } from "react-router-dom";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import axios, { AxiosError } from "axios";

const emailSchema = z.object({
  email: z
    .string()
    .email({ message: "Invalid email address" })
    .min(1, { message: "This field is required" }),
});

const PasswordChangeInitPage = () => {
  const [email, setEmail] = useState<string | null>(null);
  const { toast } = useToast();

  const emailForm = useForm<z.infer<typeof emailSchema>>({
    resolver: zodResolver(emailSchema),
    defaultValues: {
      email: "",
    },
  });

  const emailFormOnSubmit = async (values: z.infer<typeof emailSchema>) => {
    try {
      const response = await axios.post("/api/user/password-change-init", {
        email: values.email,
      });

      if (response.status === 201) setEmail(values.email);
    } catch (err) {
      toast({
        variant: "destructive",
        title: "Problem with modifyUsering",
        description: (err as AxiosError).message,
      });
    }
  };

  return (
    <div className="flex">
      <div className="flex items-center justify-center w-1/3 bg-slate-900">
        <p className="text-5xl font-bold text-white">Verify Email Page</p>
      </div>
      <div className="flex items-center justify-center h-screen w-2/3 bg-slate-950">
        {email ? (
          <span className="text-white font-semibold w-1/2">
            An email has been sent to{" "}
            <span className="text-slate-600">{email}</span>. Check your inbox
            for the verification link. If You can't see it, then check the spam
            folder.
            <br />
            <Link className="font-bold underline text-slate-600" to={"/"}>
              Back to the home page
            </Link>
          </span>
        ) : (
          <Form {...emailForm}>
            <form
              onSubmit={emailForm.handleSubmit(emailFormOnSubmit)}
              className="flex flex-col items-center gap-y-4"
            >
              <FormField
                control={emailForm.control}
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
              <Button type="submit" className="transition-all duration-300">
                Send Verification Email
              </Button>
            </form>
          </Form>
        )}
      </div>
    </div>
  );
};

export default PasswordChangeInitPage;
