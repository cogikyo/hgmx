import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bays = [el.querySelector('path[d="M9 6h6"]'), el.querySelector('path[d="M9 10h6"]')];
  animate(bays, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: stagger(0.06) });

  const light = el.querySelector('path[d="M15 14h.01"]');
  animate(light, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
