import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = ["M21 11.127", "M3.29 7", "m7.5 4.27", "M12 22V12"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const check = el.querySelector('path[d="m16 17 2 2 4-4"]');
  animate(check, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(check, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
