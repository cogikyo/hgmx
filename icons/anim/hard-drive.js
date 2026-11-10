import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const drive = el.querySelector('path[d^="M2.212"]');
  animate(drive, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const seam = el.querySelector('path[d^="M21.946"]');
  animate(seam, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });

  const lights = el.querySelectorAll('path[d$="h.01"]');
  animate(lights, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
