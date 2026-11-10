import { animate, stagger } from "../motion.js";

export default function (el) {
  const insulators = ["M3", "M7", "M17", "M21"].map((x) => el.querySelector(`path[d="${x} 3v2"]`));
  animate(insulators, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
