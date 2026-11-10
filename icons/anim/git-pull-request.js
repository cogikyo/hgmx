import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelector("line");
  animate(base, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const request = el.querySelector("path");
  animate(request, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector('circle[cx="18"]');
  animate(head, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
