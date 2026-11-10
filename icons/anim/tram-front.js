import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tram = el.querySelector("g");
  animate(tram, { scale: [0.5, 1.1, 1], opacity: [0, 1] }, { duration: τmax });

  const lights = el.querySelectorAll('path[d$="h.01"]');
  animate(lights, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
