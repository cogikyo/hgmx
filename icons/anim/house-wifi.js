import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M3 10"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const waves = ["M12 17", "M9.5 13.866", "M7 10.754"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
