import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="m12 12"]');
  animate(body, { translateX: [0, 1.5, 0] }, { duration: τmid, delay: τmin });
  const pellets = [el.querySelector('path[d="M18 12h.01"]'), el.querySelector('path[d="M22 12h.01"]')];
  animate(pellets, { translateX: [4, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
