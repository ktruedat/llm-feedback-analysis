'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth-store';
import { ProtectedRoute } from '@/components/ProtectedRoute';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { CheckCircle2, AlertCircle, LogOut, LayoutDashboard } from 'lucide-react';

const feedbackSchema = z.object({
  rating: z.number().min(1).max(5),
  comment: z.string().min(1, 'Comment is required').max(1000, 'Comment must be 1000 characters or less'),
});

type FeedbackFormData = z.infer<typeof feedbackSchema>;

function FeedbackForm() {
  const router = useRouter();
  const { clearAuth, isAdmin } = useAuthStore();
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setValue,
    watch,
  } = useForm<FeedbackFormData>({
    resolver: zodResolver(feedbackSchema),
    defaultValues: {
      rating: 5,
      comment: '',
    },
    mode: 'onChange',
  });

  const currentRating = watch('rating');

  const onSubmit = async (data: FeedbackFormData) => {
    setIsLoading(true);
    setError(null);
    setSuccess(false);

    try {
      await apiClient.createFeedback(data.rating, data.comment);
      setSuccess(true);
      // Reset form to default values
      reset({
        rating: 5,
        comment: '',
      });
      // Ensure the rating is set properly
      setValue('rating', 5);
      setTimeout(() => setSuccess(false), 5000);
    } catch (err: any) {
      if (err.response?.status === 401) {
        clearAuth();
        router.push('/login');
      } else {
        setError(err.response?.data?.message || 'Failed to submit feedback. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <CardTitle className="text-3xl">Submit Feedback</CardTitle>
              <div className="flex gap-2">
                {isAdmin && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => router.push('/admin')}
                  >
                    <LayoutDashboard className="h-4 w-4 mr-2" />
                    Admin Dashboard
                  </Button>
                )}
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    clearAuth();
                    router.push('/login');
                  }}
                >
                  <LogOut className="h-4 w-4 mr-2" />
                  Logout
                </Button>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {success && (
              <Alert className="mb-4">
                <CheckCircle2 className="h-4 w-4" />
                <AlertDescription>
                  Feedback submitted successfully! Thank you for your input.
                </AlertDescription>
              </Alert>
            )}

            {error && (
              <Alert variant="destructive" className="mb-4">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <div className="space-y-2">
                <Label>
                  Rating <span className="text-destructive">*</span>
                </Label>
                <div className="flex gap-2">
                  {[1, 2, 3, 4, 5].map((value) => {
                    const isSelected = currentRating === value;
                    return (
                      <label
                        key={value}
                        className={`flex-1 cursor-pointer text-center p-4 border-2 rounded-lg transition-all hover:bg-accent ${
                          isSelected
                            ? 'border-primary bg-primary/10'
                            : errors.rating
                            ? 'border-destructive'
                            : 'border-border'
                        }`}
                      >
                        <input
                          type="radio"
                          value={value}
                          checked={isSelected}
                          className="sr-only"
                          onChange={() => {
                            setValue('rating', value, { shouldValidate: true });
                          }}
                        />
                        <div className="text-2xl">
                          {value === 1 && '⭐'}
                          {value === 2 && '⭐⭐'}
                          {value === 3 && '⭐⭐⭐'}
                          {value === 4 && '⭐⭐⭐⭐'}
                          {value === 5 && '⭐⭐⭐⭐⭐'}
                        </div>
                        <div className="text-sm text-muted-foreground mt-1">{value}</div>
                      </label>
                    );
                  })}
                </div>
                {errors.rating && (
                  <p className="text-sm text-destructive">{errors.rating.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="comment">
                  Comment <span className="text-destructive">*</span>
                </Label>
                <Textarea
                  {...register('comment')}
                  id="comment"
                  rows={6}
                  placeholder="Please share your feedback..."
                  className="resize-none"
                />
                {errors.comment && (
                  <p className="text-sm text-destructive">{errors.comment.message}</p>
                )}
                <p className="text-xs text-muted-foreground">Maximum 1000 characters</p>
              </div>

              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? 'Submitting...' : 'Submit Feedback'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default function FeedbackPage() {
  return (
    <ProtectedRoute>
      <FeedbackForm />
    </ProtectedRoute>
  );
}
