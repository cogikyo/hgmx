import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = ["M8 18v-1", "M12 18v-6", "M16 18v-3"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
