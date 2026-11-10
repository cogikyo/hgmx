import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const rings = [el.querySelector('path[d^="M8.002"]'), el.querySelector('path[d^="M20.561"]')];
  animate(rings, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
  const flag = el.querySelector('path[d^="M12 13V2"]');
  animate(flag, { rotate: [0, -10, 6, -3, 0] }, { duration: τmid, delay: τmin });
}
