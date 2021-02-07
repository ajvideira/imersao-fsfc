import { IsNotEmpty, IsNumber, IsOptional, IsString, Min } from 'class-validator';

export class TransactionDTO {
  @IsNotEmpty()
  @IsString()
  readonly pix_key_key: string;

  @IsNotEmpty()
  @IsString()
  readonly pix_key_kind: string;

  @IsString()
  @IsOptional()
  readonly description: string;

  @IsNumber({ maxDecimalPlaces: 2 })
  @IsNotEmpty()
  @Min(0.01)
  readonly amount: number;
}
