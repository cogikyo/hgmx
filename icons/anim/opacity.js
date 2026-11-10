import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const at = (x, y) => el.querySelector(`path[d="M${x} ${y}h.01"]`);
  const dots = [at(6, 3), at(3, 6), at(12, 3), at(9, 6), at(6, 9), at(3, 12), at(18, 3), at(15, 6), at(12, 9), at(9, 12), at(6, 15), at(3, 18)];
  animate(dots, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.03) });
}
