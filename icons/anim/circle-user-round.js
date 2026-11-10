import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const outerCircle = el.querySelector('circle[r="10"]');
  animate(
    outerCircle,
    { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 115] },
    { duration: τmid },
  );

  const body = el.querySelector("path");
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelectorAll('circle[r="4"]');
  animate(head, { translateY: [1, -2, 0], opacity: [1] }, { duration: τmin, delay: τmid });
}
