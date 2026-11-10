import { animate, stagger } from "../motion.js";

export default function (el) {
  const at = (x, y) => el.querySelector(`rect[x="${x}"][y="${y}"]`);
  const tiles = [at(3, 3), at(14, 3), at(3, 14), at(14, 14)];
  animate(tiles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
