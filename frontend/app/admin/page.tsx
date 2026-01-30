'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth-store';
import { ProtectedRoute } from '@/components/ProtectedRoute';
import { Feedback } from '@/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { LogOut, LayoutDashboard, Sparkles, Trash2 } from 'lucide-react';

interface Statistics {
  total: number;
  averageRating: number;
  ratingDistribution: { [key: number]: number };
}

function AdminDashboard() {
  const router = useRouter();
  const { clearAuth } = useAuthStore();
  const [feedbacks, setFeedbacks] = useState<Feedback[]>([]);
  const [statistics, setStatistics] = useState<Statistics | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [llmAnalysis, setLlmAnalysis] = useState<{
    topics: Array<{ topic: string; count: number; examples: string[] }>;
    summary: string;
  } | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  useEffect(() => {
    loadFeedbacks();
  }, []);

  const loadFeedbacks = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.listFeedbacks(1000, 0);
      const feedbacksData = response.feedbacks || response.data?.feedbacks || [];
      setFeedbacks(feedbacksData);

      // Calculate statistics
      const total = feedbacksData.length;
      const sum = feedbacksData.reduce((acc: number, f: Feedback) => acc + f.rating, 0);
      const averageRating = total > 0 ? sum / total : 0;

      const ratingDistribution: { [key: number]: number } = { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 };
      feedbacksData.forEach((f: Feedback) => {
        ratingDistribution[f.rating] = (ratingDistribution[f.rating] || 0) + 1;
      });

      setStatistics({
        total,
        averageRating,
        ratingDistribution,
      });
    } catch (err: any) {
      if (err.response?.status === 401 || err.response?.status === 403) {
        clearAuth();
        router.push('/login');
      } else {
        setError('Failed to load feedbacks. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleAnalyze = async () => {
    setIsAnalyzing(true);
    try {
      const analysis = analyzeFeedbacks(feedbacks);
      setLlmAnalysis(analysis);
    } catch (err) {
      setError('Failed to analyze feedbacks. Please try again.');
    } finally {
      setIsAnalyzing(false);
    }
  };

  const analyzeFeedbacks = (feedbacks: Feedback[]) => {
    const topics: { [key: string]: { count: number; examples: string[] } } = {};

    feedbacks.forEach((feedback) => {
      const comment = feedback.comment.toLowerCase();
      let topic = 'General';

      if (comment.includes('service') || comment.includes('staff') || comment.includes('support')) {
        topic = 'Service Quality';
      } else if (comment.includes('product') || comment.includes('feature')) {
        topic = 'Product Features';
      } else if (comment.includes('price') || comment.includes('cost') || comment.includes('expensive')) {
        topic = 'Pricing';
      } else if (comment.includes('bug') || comment.includes('error') || comment.includes('issue')) {
        topic = 'Technical Issues';
      } else if (comment.includes('fast') || comment.includes('slow') || comment.includes('speed')) {
        topic = 'Performance';
      }

      if (!topics[topic]) {
        topics[topic] = { count: 0, examples: [] };
      }
      topics[topic].count++;
      if (topics[topic].examples.length < 3) {
        topics[topic].examples.push(feedback.comment.substring(0, 100));
      }
    });

    const topicArray = Object.entries(topics).map(([topic, data]) => ({
      topic,
      count: data.count,
      examples: data.examples,
    }));

    const positiveCount = feedbacks.filter((f) => f.rating >= 4).length;
    const negativeCount = feedbacks.filter((f) => f.rating <= 2).length;
    const neutralCount = feedbacks.length - positiveCount - negativeCount;

    const summary = `Analysis of ${feedbacks.length} feedback entries:
- ${positiveCount} positive feedbacks (4-5 stars)
- ${neutralCount} neutral feedbacks (3 stars)
- ${negativeCount} negative feedbacks (1-2 stars)
- Average rating: ${statistics?.averageRating.toFixed(2) || '0.00'}/5.00

Key topics identified: ${topicArray.map((t) => t.topic).join(', ')}`;

    return {
      topics: topicArray,
      summary,
    };
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this feedback?')) {
      return;
    }

    try {
      await apiClient.deleteFeedback(id);
      await loadFeedbacks();
    } catch (err: any) {
      if (err.response?.status === 403) {
        setError('You do not have permission to delete feedbacks.');
      } else {
        setError('Failed to delete feedback. Please try again.');
      }
    }
  };

  if (isLoading) {
    return (
      <ProtectedRoute>
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            <h1 className="text-2xl font-bold mb-4">Loading...</h1>
          </div>
        </div>
      </ProtectedRoute>
    );
  }

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto space-y-6">
          <div className="flex justify-between items-center">
            <h1 className="text-3xl font-bold">Admin Dashboard</h1>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => router.push('/feedback')}>
                Submit Feedback
              </Button>
              <Button variant="outline" onClick={() => { clearAuth(); router.push('/login'); }}>
                <LogOut className="h-4 w-4 mr-2" />
                Logout
              </Button>
            </div>
          </div>

          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          {/* Statistics */}
          {statistics && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle>Total Feedbacks</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-3xl font-bold">{statistics.total}</p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Average Rating</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-3xl font-bold">{statistics.averageRating.toFixed(2)}/5.00</p>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Rating Distribution</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-2">
                    {[5, 4, 3, 2, 1].map((rating) => (
                      <div key={rating} className="flex items-center gap-2">
                        <span className="text-sm w-4">{rating}</span>
                        <div className="flex-1 bg-secondary rounded-full h-2">
                          <div
                            className="bg-primary h-2 rounded-full"
                            style={{
                              width: `${
                                statistics.total > 0
                                  ? (statistics.ratingDistribution[rating] / statistics.total) * 100
                                  : 0
                              }%`,
                            }}
                          />
                        </div>
                        <span className="text-sm w-8 text-right">
                          {statistics.ratingDistribution[rating]}
                        </span>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}

          {/* AI Analysis */}
          <Card>
            <CardHeader>
              <div className="flex justify-between items-center">
                <div>
                  <CardTitle>AI Analysis</CardTitle>
                  <CardDescription>Analyze feedback patterns and topics</CardDescription>
                </div>
                <Button
                  onClick={handleAnalyze}
                  disabled={isAnalyzing || feedbacks.length === 0}
                >
                  <Sparkles className="h-4 w-4 mr-2" />
                  {isAnalyzing ? 'Analyzing...' : 'Run Analysis'}
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              {llmAnalysis && (
                <div className="space-y-6">
                  <div>
                    <h3 className="text-lg font-semibold mb-2">Summary</h3>
                    <div className="bg-muted rounded-md p-4 whitespace-pre-line text-sm">
                      {llmAnalysis.summary}
                    </div>
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold mb-2">Topic Clustering</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      {llmAnalysis.topics.map((topic, index) => (
                        <Card key={index}>
                          <CardHeader>
                            <CardTitle className="text-base">
                              {topic.topic} <Badge>{topic.count}</Badge>
                            </CardTitle>
                          </CardHeader>
                          <CardContent>
                            <ul className="text-sm space-y-1">
                              {topic.examples.map((example, i) => (
                                <li key={i} className="italic text-muted-foreground">
                                  "{example}..."
                                </li>
                              ))}
                            </ul>
                          </CardContent>
                        </Card>
                      ))}
                    </div>
                  </div>
                </div>
              )}
              {!llmAnalysis && feedbacks.length > 0 && (
                <p className="text-muted-foreground text-sm">Click "Run Analysis" to analyze feedback patterns.</p>
              )}
              {feedbacks.length === 0 && (
                <p className="text-muted-foreground text-sm">No feedbacks available for analysis.</p>
              )}
            </CardContent>
          </Card>

          {/* Recent Submissions */}
          <Card>
            <CardHeader>
              <CardTitle>Recent Submissions</CardTitle>
            </CardHeader>
            <CardContent>
              {feedbacks.length === 0 ? (
                <p className="text-muted-foreground">No feedbacks submitted yet.</p>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Rating</TableHead>
                      <TableHead>Comment</TableHead>
                      <TableHead>Date</TableHead>
                      <TableHead>Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {feedbacks.slice(0, 20).map((feedback) => (
                      <TableRow key={feedback.id}>
                        <TableCell>
                          <div className="flex items-center gap-2">
                            <span className="text-lg">{'⭐'.repeat(feedback.rating)}</span>
                            <Badge variant="outline">{feedback.rating}</Badge>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="max-w-md truncate">{feedback.comment}</div>
                        </TableCell>
                        <TableCell>
                          {new Date(feedback.created_at).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDelete(feedback.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </ProtectedRoute>
  );
}

export default function AdminPage() {
  return <AdminDashboard />;
}
