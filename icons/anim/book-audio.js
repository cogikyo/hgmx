import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = ["M8 8v3", "M12 6v7", "M16 8v3"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleY: [1, 0.6, 1.2, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
