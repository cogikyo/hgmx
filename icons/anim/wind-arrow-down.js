import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const gusts = ["M17.5", "M12.8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(gusts, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
  const arrow = el.querySelectorAll('path[d^="M10 2"], path[d^="m6 6"]');
  animate(arrow, { translateY: [-4, 2, 0], opacity: [0, 1] }, { duration: τmax });
}
