import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import axios from "axios";

const registerFormSchema = z
  .object({
    username: z
      .string()
      .min(5, { message: "Must be at least 5 characters long." })
      .max(50, { message: "Can't be longer than 50 characters." }),
    email: z
      .string()
      .email({ message: "Ivalid email." })
      .min(1, { message: "This field is required." }),
    password: z.string().min(1, { message: "This field is required." }),
    repeatPass: z.string(),
  })
  .refine((data) => data.password === data.repeatPass, {
    message: "Passwords don't match",
    path: ["repeatPass"],
  });

const RegisterPage = () => {
  const registerForm = useForm<z.infer<typeof registerFormSchema>>({
    resolver: zodResolver(registerFormSchema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
      repeatPass: "",
    },
  });

  const registerFormOnSubmit = (values: z.infer<typeof registerFormSchema>) => {
    axios
      .post("/api/register", {
        body: {
          username: values.username,
          email: values.email,
          password: values.password,
        },
      })
      .then((response) => {
        console.log(response);
      });
  };

  return (
    <div className="flex">
      <div className="flex items-center justify-center w-1/3 bg-slate-900">
        <p className="text-5xl font-bold text-white">Register Page</p>
      </div>
      <div className="flex items-center justify-center h-screen w-2/3 bg-slate-950">
        <Form {...registerForm}>
          <form
            onSubmit={registerForm.handleSubmit(registerFormOnSubmit)}
            className="flex flex-col items-center gap-y-4"
          >
            <FormField
              control={registerForm.control}
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
              control={registerForm.control}
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
              control={registerForm.control}
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
              control={registerForm.control}
              name="repeatPass"
              render={({ field }) => (
                <FormItem className="transition-all duration-300">
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="Repeat password"
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
              Submit
            </Button>
          </form>
        </Form>
      </div>
    </div>
  );
};

export default RegisterPage;
