import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const ground = el.querySelector('path[d="M22 21H2"]');
  animate(ground, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const steps = ["M19 17h", "M12 17v", "M19 9h"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(steps, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
