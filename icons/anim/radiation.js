import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const blades = el.querySelectorAll("path[transform-origin]");
  animate(blades, { scale: [1, 1.15, 0.95, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const core = el.querySelector('path[d^="M12 12"]');
  animate(core, { scale: [1, 1.75, 1] }, { duration: τmin, delay: τmid });
}
