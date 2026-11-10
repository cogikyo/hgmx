import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector('path[d^="M12 2a"]');
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const hands = el.querySelector('path[d^="M12 6v6"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });
  const trail = ["M8.644", "M2.83", "M2.5", "M4.636"].map((d) => el.querySelector(`path[d^="${d} "]`));
  animate(trail, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.03, { startDelay: τmid }) });
}
