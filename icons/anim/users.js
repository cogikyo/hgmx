import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const secondUserHead = el.querySelector('path[d*="16 3.128"]');
  animate(
    secondUserHead,
    { pathLength: [0, 1], opacity: [0, 1], translateX: [3, 0] },
    { duration: τmid },
  );

  const secondUserBody = el.querySelector('path[d*="22 21v-2"]');
  animate(
    secondUserBody,
    { pathLength: [0, 1], opacity: [0, 1], translateX: [3, 0] },
    { duration: τmid },
  );

  const mainBody = el.querySelector('path[d*="16 21v-2"]');
  animate(mainBody, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const mainHead = el.querySelector("circle");
  animate(mainHead, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
