import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M22 17"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dots = ["M8 11", "M12 11", "M16 11"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(dots, { translateY: [0, -2, 1, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
