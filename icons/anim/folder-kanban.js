import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M4 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const columns = ["M8 10v4", "M12 10v2", "M16 10v6"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(columns, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
