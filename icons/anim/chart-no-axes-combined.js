import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bars = ["M4 18.463V21", "M8 14.656V21", "M12 16v5", "M16 14.639V21", "M20 10.656V21"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(bars, { scaleY: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const line = el.querySelector('path[d^="m22 3"]');
  animate(line, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
