import { z } from 'zod';

// A matrix: a non-empty array of non-empty numeric rows, all the same length.
const matrixData = z
  .array(z.array(z.number()).min(1))
  .min(1)
  .refine((rows) => rows.every((row) => row.length === rows[0].length), {
    message: 'all matrix rows must have the same length',
  });

export const statsRequestSchema = z.object({
  matrices: z
    .array(
      z.object({
        name: z.string().min(1),
        data: matrixData,
      }),
    )
    .min(1),
});

export type StatsRequest = z.infer<typeof statsRequestSchema>;
