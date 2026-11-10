import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const wand = el.querySelector('path[d^="m3 21"]');
  animate(wand, { rotate: [0, -10, 4, 0] }, { duration: τmax });

  const glints = ["M15 9h", "M15 4", "M17.8 6.2", "M20 9", "M17.8 11.8", "M15 16", "M8 9", "M12.2"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(glints, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
