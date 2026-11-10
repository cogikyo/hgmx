import { animate, stagger } from "../motion.js";

export default function (el) {
  const waves = [el.querySelector('path[d^="M17 15"]'), el.querySelector('path[d^="M19.8"]')];
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
