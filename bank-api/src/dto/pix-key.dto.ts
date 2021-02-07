import { IsIn, IsNotEmpty, IsString } from 'class-validator';
import { PixKeyKind } from 'src/models/pix-key.model';

export class PixKeyDTO {
  @IsString()
  @IsNotEmpty()
  readonly key: string;

  @IsString()
  @IsNotEmpty()
  @IsIn(Object.values(PixKeyKind))
  readonly kind: PixKeyKind;
}
