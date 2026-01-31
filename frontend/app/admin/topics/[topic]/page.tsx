'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth-store';
import { AdminRoute } from '@/components/AdminRoute';
import { TopicDetails } from '@/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { ArrowLeft, Loader2 } from 'lucide-react';

function TopicDetailPage() {
  const router = useRouter();
  const params = useParams();
  const { clearAuth } = useAuthStore();
  const topicEnum = params.topic as string;
  const [topicDetails, setTopicDetails] = useState<TopicDetails | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (topicEnum) {
      loadTopicDetails();
    }
  }, [topicEnum]);

  const loadTopicDetails = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.getTopicDetails(topicEnum);
      const detailsData = response?.data || response;
      if (detailsData) {
        setTopicDetails(detailsData);
      }
    } catch (err: any) {
      if (err.response?.status === 401 || err.response?.status === 403) {
        clearAuth();
        router.push('/login');
      } else if (err.response?.status === 404) {
        setError('Topic not found in the latest analysis.');
      } else {
        setError('Failed to load topic details. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return (
      <AdminRoute>
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4" />
            <h1 className="text-2xl font-bold mb-4">Loading topic details...</h1>
          </div>
        </div>
      </AdminRoute>
    );
  }

  if (error) {
    return (
      <AdminRoute>
        <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
          <div className="max-w-7xl mx-auto">
            <Button variant="outline" onClick={() => router.back()} className="mb-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          </div>
        </div>
      </AdminRoute>
    );
  }

  if (!topicDetails) {
    return (
      <AdminRoute>
        <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
          <div className="max-w-7xl mx-auto">
            <Button variant="outline" onClick={() => router.back()} className="mb-4">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
            <p className="text-muted-foreground">No topic details found.</p>
          </div>
        </div>
      </AdminRoute>
    );
  }

  return (
    <AdminRoute>
      <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="flex items-center gap-4">
            <Button variant="outline" onClick={() => router.back()}>
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back
            </Button>
            <h1 className="text-3xl font-bold">{topicDetails.topic_name}</h1>
          </div>

          {/* Topic Overview */}
          <Card>
            <CardHeader>
              <CardTitle>Topic Overview</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h3 className="text-sm font-semibold text-muted-foreground mb-2">Description</h3>
                <div className="bg-muted rounded-md p-4 whitespace-pre-line text-sm">
                  {topicDetails.topic_description}
                </div>
              </div>
              <div className="flex gap-4 flex-wrap">
                <div>
                  <span className="text-sm text-muted-foreground">Feedback Count: </span>
                  <span className="font-semibold">{topicDetails.feedback_count}</span>
                </div>
                <div>
                  <span className="text-sm text-muted-foreground">Average Rating: </span>
                  <span className="font-semibold">{topicDetails.average_rating.toFixed(2)}/5.00</span>
                </div>
                <div>
                  <span className="text-sm text-muted-foreground">Sentiment: </span>
                  <Badge
                    variant={
                      topicDetails.sentiment === 'positive'
                        ? 'default'
                        : topicDetails.sentiment === 'negative'
                          ? 'destructive'
                          : 'secondary'
                    }
                  >
                    {topicDetails.sentiment}
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Topic Summary */}
          {topicDetails.summary && (
            <Card>
              <CardHeader>
                <CardTitle>AI Analysis Summary</CardTitle>
                <CardDescription>LLM-generated summary for this topic</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="bg-muted rounded-md p-4 whitespace-pre-line text-sm">
                  {topicDetails.summary}
                </div>
              </CardContent>
            </Card>
          )}

          {/* Feedbacks */}
          <Card>
            <CardHeader>
              <CardTitle>Associated Feedbacks</CardTitle>
              <CardDescription>
                {topicDetails.feedbacks.length} feedback{topicDetails.feedbacks.length !== 1 ? 's' : ''} associated with this topic
              </CardDescription>
            </CardHeader>
            <CardContent>
              {topicDetails.feedbacks.length === 0 ? (
                <p className="text-muted-foreground">No feedbacks associated with this topic.</p>
              ) : (
                <div className="space-y-4">
                  {topicDetails.feedbacks.map((feedback) => (
                    <Card key={feedback.id} className="bg-white">
                      <CardContent className="pt-6">
                        <div className="flex items-start justify-between gap-4">
                          <div className="flex-1">
                            <div className="flex items-center gap-2 mb-2">
                              <span className="text-lg">{'⭐'.repeat(feedback.rating)}</span>
                              <Badge variant="outline">{feedback.rating}/5</Badge>
                              <span className="text-sm text-muted-foreground">
                                {new Date(feedback.created_at).toLocaleDateString()}
                              </span>
                            </div>
                            <p className="text-sm whitespace-pre-line">{feedback.comment}</p>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </AdminRoute>
  );
}

export default function TopicPage() {
  return <TopicDetailPage />;
}
