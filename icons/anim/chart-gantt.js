import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tasks = ["M10 6h8", "M8 11h7", "M12 16h6"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(tasks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.09, { startDelay: τmin }) });
}
