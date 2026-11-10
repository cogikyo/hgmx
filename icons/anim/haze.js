import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const rays = ["M2 13", "m5.2", "M12 5", "m17.4", "M20 13"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
  const haze = [el.querySelector('path[d^="M22 17"]'), el.querySelector('path[d^="M22 21"]')];
  animate(haze, { pathLength: [0, 1], opacity: [0, 1], translateX: [3, 0] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
