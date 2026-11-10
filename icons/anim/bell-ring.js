import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bell = el.querySelector("g");
  animate(bell, { rotate: [0, -18, 14, -8, 0] }, { duration: τmax });
  const clapper = el.querySelector('path[d^="M10.268"]');
  animate(clapper, { translateX: [0, 2, -2, 0] }, { duration: τmid, delay: τmin });
  const waves = [el.querySelector('path[d^="M22 8"]'), el.querySelector('path[d^="M4 2"]')];
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
